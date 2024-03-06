# build
# go build
set -x
#download
#hu

m3u8="https://yzzy.play-cdn17.com/20230731/13995_865a8eba/index.m3u8"
m3u8="https://web-vod-xdrive.xunlei.com/ts_downloader?client_id=Xqp0kJBXWhwaTpB6&url=https%3A%2F%2Fvod1491-aliyun06-vip-lixian.xunlei.com%2Fdownload%2F%3Ffid%3D7oXbHKtZGHtQuDB3Txy-YEzG*5DM4a8aAAAAAAUdajEcWhaXrFHiiFlfmUyjmA8M%26mid%3D666%26threshold%3D251%26tid%3D73A9A76F30B525C81F9DA55AEFCAC1DE%26srcid%3D0%26verno%3D2%26pk%3Dxdrive%26e%3D1709723845%26g%3D051D6A311C5A1697AC51E288595F994CA3980F0C%26i%3DEE85DB1CAB59187B50B830774F1CBF604CC6FB90%26ui%3D891934429%26t%3D1%26hy%3D0%26ms%3D1670490%26th%3D167049%26pt%3D0%26f%3D447734220%26alt%3D0%26pks%3D654%26rts%3D%26fileid%3DVNs6IZGL568HZ5ILWyRKA68_A1%26spr%3Dplaytrans%26vip%3DFREE%26cliplugver%3D%26tras%3D1%26vc%3Dhevc%26source%3Dxdrive%26clientid%3DXqp0kJBXWhwaTpB6%26projectid%3D2rvk4e3gkdnl7u1kl0k%26userid%3D891934429%26clientver%3D%26fext%3Dmp4%26rg%3D0-16704900%26at%3DECFEA822ABA538DA7353ABD11732CB36#zc.m3u8"


ff="-vf scale=-1:480 -to 00:00:10 -c:v libx264 -acodec copy -preset veryslow -crf 28"
ff="-vf scale=-1:480 -c:v libx264 -preset veryslow -crf 28"

now=`date "+%Y-%m-%d %H:%M:%S"`
echo "$now start" >> run.log
merged=zhouChuChuSanHai.mp4

#./goose -u $m3u8 -ff "$ff" -n $merged -threads 2 -t 0 >> run.log
./goose -prefix " " -u $m3u8 -ff "$ff" -n $merged -threads 2 -t 0 >> run.log


now=`date "+%Y-%m-%d %H:%M:%S"`
echo "$now end" >> run.log
#ffmpeg -ss 00:47:40 -to 00:52:38  -i "yanzilong.ts" -s 640x360 -crf 28  -acodec copy  yanzilong07.mp4
#ffmpeg -i huxing.ts -vf scale=-1:360 -acodec copy -preset veryslow -crf 28 huxing.mp4
