package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// 下载m3u8文件
func downloadM3U8(url string) (string, error) {
	baseName := filepath.Base(url)
	dir := strings.TrimSuffix(baseName, ".m3u8")
	err := os.Mkdir(dir, 0755)
	if err != nil && !os.IsExist(err) {
		return "", err
	}
	err = os.Chdir(dir)
	if err != nil {
		return "", err
	}
	cmd := exec.Command("wget", url)
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	return baseName, nil
}

// 获取ts下载文件
func getUrls(prefix, filename, suffix string) ([]byte, []string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, nil
	}

	keyUrlRE := regexp.MustCompile(`AES-128,URI="(.*?)"`)
	keyUrl := keyUrlRE.FindSubmatch(data)
	fmt.Printf("key:%s\n\n", keyUrl)
	if len(keyUrl) < 2 {
	    fmt.Printf("not encrypted\n")
		//return nil, nil, errors.New("failed to match key url")
	}

	// 解析url
	tsUrlRE := regexp.MustCompile(`.+ts\?.+`)
	tsUrls := tsUrlRE.FindAllSubmatch(data, -1)
	if len(tsUrls) == 0 {
		return nil, nil, errors.New("failed to match ts urls")
	}

	urls := make([]string, len(tsUrls))
	for i, lv1 := range tsUrls {
		urls[i] = prefix + string(lv1[0]) + suffix
	}

    if len(keyUrl) < 2 {
	    return nil, urls, nil
	}

	var ku = string(keyUrl[1])
	ku = ku + "&uid=u_63804e2422b27_a6FrFVVpgP"
	fmt.Printf("keyUrl:%s\n", ku)
	cmd := exec.Command("wget", ku, "-O", "key")
	err = cmd.Run()
	if err != nil {
		return nil, nil, err
	}

	key, err := ioutil.ReadFile("./key")
	if err != nil {
		return nil, nil, err
	}

	return key, urls, nil
}

type task struct {
	num int
	url string
}

type result struct {
	num      int
	url      string
	filename string
	err      error
}

func getFileSize(path string) int64{
	fi,err := os.Stat(path);
	if err != nil {
		return 0;
	}

	return fi.Size();
}

func downloadChunks(key []byte, urls []string) (int, error) {
	downloadCh := make(chan task, len(urls))
	resultCh := make(chan result, len(urls))
	var wg sync.WaitGroup
	wg.Add(len(urls))

	// 发出下载任务
	for i, url := range urls {
		downloadCh <- task{
			num: i,
			url: url,
		}
	}
	close(downloadCh)

	// 启动10个协程下载文件
	fmt.Println("total: ", len(urls))
	for i := 0; i < 10; i++ {
		go func() {
			for t := range downloadCh {
				filename := strconv.Itoa(t.num) + ".ts"
				fileSize := getFileSize(filename)
				var err error
				if fileSize > 1024 {
					//fmt.Println("skipped ", filename)
				}else{
					// 下载文件
					fmt.Printf("downloading %d %s\n", t.num, t.url)
					cmd := exec.Command("wget", t.url, "-O", filename)
					err = cmd.Run()
					if err != nil {
						fmt.Println("download error: ", err)
					}else{
						fmt.Printf("downloaded %d\n", t.num)
					}
				}
				resultCh <- result{
					num:      t.num,
					url:      t.url,
					filename: filename,
					err:      err,
				}
				wg.Done()
			}
		}()
	}

	wg.Wait()

	// 检查结果
	fmt.Printf("tried: %d\n", len(resultCh))
	for i := 0; i < len(urls); i++ {
		rs := <-resultCh
		if rs.err != nil {
			return 0, fmt.Errorf("failed to download %s: %v", rs.url, rs.err)
		}

		if key == nil {
			continue
		}

		// 解密
		block, err := aes.NewCipher(key)
		if err != nil {
			return 0, err
		}

		fileData, err := ioutil.ReadFile(rs.filename)
		if err != nil {
			return 0, err
		}
		pt := make([]byte, len(fileData))

		bm := cipher.NewCBCDecrypter(block, bytes.Repeat([]byte{0}, 16))
		bm.CryptBlocks(pt, fileData)

		err = ioutil.WriteFile(rs.filename, pt, 0755)
		if err != nil {
			return 0, err
		}
	}

	return len(urls), nil
}

func mergeFile(count int) error {
	// ffmpeg -i "concat:ttt.ts|tt2.ts" -c copy output.ts
	files := make([]string, count)
	for i := range files {
		files[i] = strconv.Itoa(i) + ".ts"
	}

	// Too many open files
	// ulimit -n 1024

	fmt.Println("ffmpeg", "-i", fmt.Sprintf("\"concat:%s\"", strings.Join(files, "|")), "merge.ts")

	cmd := exec.Command("ffmpeg", "-i", fmt.Sprintf("concat:%s", strings.Join(files, "|")), "merge.ts")
	o, e := cmd.CombinedOutput()
	fmt.Println(string(o))
	return e
}

func getPrefix(url string) string {
	i := strings.LastIndex(url, "/")
	return url[:i+1]
}

func main() {
	var url string
	var newName string
	var downloadPrefix string
	var downloadSuffix string
	var test int
	var ks string
	flag.StringVar(&url, "u", "", "m3u8 url")
	flag.StringVar(&newName, "n", "", "new name")
	flag.StringVar(&downloadPrefix, "prefix", "", "prefix of download url")
	flag.StringVar(&downloadSuffix, "suffix", "", "suffix of download url")
	flag.IntVar(&test, "t", 0, "test (only download n chunk)")
	flag.StringVar(&ks, "key", "", "key")
	flag.Parse()

        found, err := regexp.MatchString("m3u8($|\\?.*)", url)
        if err != nil {
                panic(err)
        }
        if !found {
                fmt.Println("please enter valid m3u8 url")
                return
        }

	// 1. 下载m3u8
	filename, err := downloadM3U8(url)
	if err != nil {
		panic(err)
	}

	var prefix = getPrefix(url)
	if len(downloadPrefix) > 0 {
		prefix=downloadPrefix
	}
	//prefix = "https://c-vod.hw-cdn.xiaoeknow.com/9764a7a5vodtransgzp1252524126/a77e2e61387702296415479954/drm/"
	// 2. 解析出key和分片url
	key, tsUrls, err := getUrls(prefix, filename, downloadSuffix)
	if err != nil {
		panic(err)
	}

	if len(ks) > 0 {
		data, err := base64.StdEncoding.DecodeString(ks)
		if err != nil {
			fmt.Println("key error:", err)
			return
		}

		key = data
	}

	if test > 0 {
		fmt.Println("\n==== test ======\n");
		tsUrls = tsUrls[:test]
	}else{
		fmt.Println("\n==== not test ======\n");
	}

	// 3. 并发下载分片文件并解密
	count, err := downloadChunks(key, tsUrls)
	if err != nil {
		panic(err)
	}

	// 4. 合并文件
	err = mergeFile(count)
	fmt.Println(err)

	if err == nil {
		if len(newName) > 0 {
			err = os.Rename("merge.ts", "../"+newName)
		} else {
			err = os.Rename("merge.ts", "../merge.ts")
		}
		fmt.Println("move", err)
		if err == nil {
			os.Chdir("../")
			err = os.RemoveAll(strings.TrimSuffix(filepath.Base(url), ".m3u8"))
			fmt.Println("remove", err)
		}
	}
}

// 
//ffmpeg -i "concat:0.ts|1.ts|2.ts|3.ts|4.ts"   -s 640x360  -acodec copy -preset veryslow -crf 28  merge.ts
//ffmpeg -i "concat:0.ts|1.ts|2.ts|3.ts|4.ts"   merge.ts
//btoa(String.fromCharCode.apply(null, new Uint8Array(a.decryptdata.key.buffer)))
