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
	"bufio"
	"sync"
)

// 下载m3u8文件
func downloadM3U8(url string) (string, error) {
	baseName := filepath.Base(url)
	baseName = "download.m3u8"
        fmt.Println(baseName)
	dir := strings.TrimSuffix(baseName, ".m3u8")
	if len(dir) > 20 {
 	   dir = dir[:20]
	}
	err := os.Mkdir(dir, 0755)
	if err != nil && !os.IsExist(err) {
		return "", err
	}
	err = os.Chdir(dir)
	if err != nil {
		return "", err
	}

	userAgent:="Mozilla/5.0 (Linux; Android 8.0.0; SM-G955U Build/R16NW) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Mobile Safari/537.36"
	cmd := exec.Command("wget","-U",userAgent,"-O", baseName, url)
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
		fmt.Printf("read file error:%s\n", filename)
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
	tsUrlRE := regexp.MustCompile(`.+\?.+`)
	tsUrls := tsUrlRE.FindAllSubmatch(data, -1)
	if len(tsUrls) == 0 {
		return nil, nil, errors.New("failed to match ts urls")
	}

	urls := make([]string, len(tsUrls))
	for i, lv1 := range tsUrls {
		if strings.HasPrefix(string(lv1[0]), "http") {
			urls[i] = string(lv1[0]) + suffix
	     	}else{
			
			urls[i] = prefix + string(lv1[0]) + suffix
		}
	}
	fmt.Printf("urls:%d\n", len(urls))

    	if len(keyUrl) < 2 {
	    return nil, urls, nil
	}

	var ku = string(keyUrl[1])
	ku = ku + "&uid=u_63804e2422b27_a6FrFVVpgP"
	fmt.Printf("keyUrl:%s\n", ku)

	userAgent:="Mozilla/5.0 (Linux; Android 8.0.0; SM-G955U Build/R16NW) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Mobile Safari/537.36"
	cmd := exec.Command("wget","-U",userAgent, ku, "-O", "key")
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

func downloadChunks(threads int, key []byte, urls []string) (int, error) {
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
	for i := 0; i < threads; i++ {
		go func() {
			for t := range downloadCh {
				filename := strconv.Itoa(t.num) + ".ts"
				fileSize := getFileSize(filename)
				var err error
				if fileSize > 1024*1024*1024 {
					fmt.Printf("skipped %d \n", t.num)
				}else{
					// 下载文件
					fmt.Printf("downloading %d %s\n", t.num, t.url)

					userAgent:="Mozilla/5.0 (Linux; Android 8.0.0; SM-G955U Build/R16NW) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Mobile Safari/537.36"
					cmd := exec.Command("wget","-c", "-U",userAgent, t.url, "-O", filename)

					err = cmd.Run()
					if err != nil {
						fmt.Printf("download %d error: %s\n", t.num, err)
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

func mergeFile(count int, params string) error {
	// ffmpeg -i "concat:ttt.ts|tt2.ts" -c copy output.ts
	files := make([]string, count)
	for i := range files {
		files[i] = strconv.Itoa(i) + ".ts"
	}

	// Too many open files
	// ulimit -n 1024
	input := fmt.Sprintf("concat:%s", strings.Join(files, "|"))
	/*
	args := []string{"-i", input}
	pattern := "[ ]+"
	re := regexp.MustCompile(pattern)
	pargs := re.Split(params, -1)
	fmt.Printf("分割后的数组为:%d个:%v\n",len(pargs), pargs)
	args = append(args, pargs...)
	args = append(args, "merge.ts")
	*/
	shell := fmt.Sprintf("cd /home/liqinghu/git/goose/download && ffmpeg -i \"%s\" %s %s", input, params, "merge.mp4")
	err := ioutil.WriteFile("/tmp/temp_merge.sh", []byte(shell), 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(shell)
	fmt.Println("bash /tmp/temp_merge.sh")
	cmd := exec.Command("bash", "/tmp/temp_merge.sh")
	//cmd := exec.Command("ffmpeg", "-i", input, "merge.ts")


	// 获取标准输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("获取标准输出失败：", err)
		return err
	}
	
	// 读取输出流
	reader := bufio.NewReader(stdout)
	go func() {
		for {
			r, _, err := reader.ReadRune()
			if err != nil {
				fmt.Println("读取标准输出失败：", err)
				return
			}
			fmt.Printf("%s", string(r))
		}
	}()

    // 获取标准错误输出
    stderr, err1 := cmd.StderrPipe()
    if err1 != nil {
        fmt.Println("获取错误输出失败：", err)
        return err1
    }

	// 读取输出流
	readerE := bufio.NewReader(stderr)
	go func() {
		for {
			r, _, err := readerE.ReadRune()
			if err != nil {
				fmt.Println("读取错误输出失败：", err)
				return
			}
			fmt.Printf("%s", string(r))
		}
	}()


	// 开始执行命令
	err = cmd.Start()
	if err != nil {
		fmt.Println("执行命令失败：", err)
		return err
	}
	
	// 等待命令执行完成
	err = cmd.Wait()
	if err != nil {
		fmt.Println("命令执行出错：", err)
		
	}

	return err

	/*
	o, e := cmd.CombinedOutput()
	fmt.Println(string(o))
	return e
	*/
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
	var params string
	var test int
	var threads int
	var ks string
	flag.StringVar(&url, "u", "", "m3u8 url")
	flag.StringVar(&newName, "n", "", "new name")
	flag.StringVar(&downloadPrefix, "prefix", "", "prefix of download url")
	flag.StringVar(&downloadSuffix, "suffix", "", "suffix of download url")
	flag.IntVar(&test, "t", 0, "test (only download n chunk)")
	flag.StringVar(&ks, "key", "", "key")
	flag.IntVar(&threads, "threads", 1, "threads count")
	flag.StringVar(&params, "ff", "", "ffmpeg params")
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
	count, err := downloadChunks(threads, key, tsUrls)
	if err != nil {
		panic(err)
	}

	// 4. 合并文件
	err = mergeFile(count, params)
	fmt.Println(err)

	if err == nil {
		if len(newName) > 0 {
			err = os.Rename("merge.mp4", "../"+newName)
		} else {
			err = os.Rename("merge.mp4", "../merge.mp4")
		}
		fmt.Println("move", err)
		/*
		if err == nil {
			os.Chdir("../")
			err = os.RemoveAll(strings.TrimSuffix(filepath.Base(url), ".m3u8"))
			fmt.Println("remove", err)
		}
		*/
	}
}

// 
//ffmpeg -i "concat:0.ts|1.ts|2.ts|3.ts|4.ts"   -s 640x360  -acodec copy -preset veryslow -crf 28  merge.ts
//ffmpeg -i "concat:0.ts|1.ts|2.ts|3.ts|4.ts"   merge.ts
//btoa(String.fromCharCode.apply(null, new Uint8Array(a.decryptdata.key.buffer)))
