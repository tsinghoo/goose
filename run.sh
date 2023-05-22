# build
# go build
set -x
#download
m3u8="https://video.qianliao.net/43842e87vodtranscq1500005495/cd65d8203270835008906974310/video_950555_1.m3u8?sign=61d9c9cc6832f902720b856ae8ea2eee&t=645f3a4c"
m3u8="https://mps-trans.yzcdn.cn/multi_trans_hls_hd/67j57BpSVwY-VWZDvLklgv96sG8bcS8_3Lymog_HD.m3u8?sign=257ace0911d3831f37a5db688b610034&t=6466d928"
m3u8="https://video.qianliao.net/43842e87vodtranscq1500005495/db125433243791581751601387/video_950555_1.m3u8?sign=99dea2a2fad03d5e6900b1c9539d154c&t=6466440a"

merged=temp.ts
ff="-vf scale=-1:360 -c:v libx264 -acodec copy -preset veryslow -crf 28"

./xiaoetong -u $m3u8 -ff "$ff" -n $merged -threads 2 -t 0

#ffmpeg -ss 00:47:40 -to 00:52:38  -i "yanzilong.ts" -s 640x360 -crf 28  -acodec copy  yanzilong07.mp4
#ffmpeg -i huxing.ts -vf scale=-1:360 -acodec copy -preset veryslow -crf 28 huxing.mp4
