# build
# go build
set -x
#download
#hu
m3u8="https://video.qianliao.net/43842e87vodtranscq1500005495/0facd5a2387702294794384097/video_950555_0.m3u8?sign=3a6ab0769b7a44ddf9f6d6724eb0a051&t=65299b0c"
m3u8="https://video.qianliao.net/43842e87vodtranscq1500005495/b70962d4387702299113206595/video_950555_0.m3u8?sign=b6aba39fe4905a845c90784b2a0d9458&t=6529a557"
m3u8="https://video.qianliao.net/43842e87vodtranscq1500005495/e520d3e9387702297491947759/video_950555_0.m3u8?sign=7f88868c16ceac708de846f6111d828f&t=6529a7dc"
m3u8="https://video2.qianliao.net/43a5a542vodtranscq1500020547/d76f400b5576678020618140668/video_1440284_0.m3u8?sign=363f8d0f3a42cbbacc00c8447d35164b&t=652a3ec8"
#video_950555_0.m3u8?sign=3a6ab0769b7a44ddf9f6d6724eb0a051&t=65299b0c

ff="-vf scale=-1:480 -to 00:00:10 -c:v libx264 -acodec copy -preset veryslow -crf 28"
ff="-vf scale=-1:480 -c:v libx264 -preset veryslow -crf 28"

now=`date "+%Y-%m-%d %H:%M:%S"`
echo "$now start" >> run.log
merged=qixue4.mp4

./goose -u $m3u8 -ff "$ff" -n $merged -threads 2 -t 0 >> run.log


now=`date "+%Y-%m-%d %H:%M:%S"`
echo "$now end" >> run.log
#ffmpeg -ss 00:47:40 -to 00:52:38  -i "yanzilong.ts" -s 640x360 -crf 28  -acodec copy  yanzilong07.mp4
#ffmpeg -i huxing.ts -vf scale=-1:360 -acodec copy -preset veryslow -crf 28 huxing.mp4
