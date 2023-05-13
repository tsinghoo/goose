# build
# go build
set -x
#download
m3u8='https://video.qianliao.net/43842e87vodtranscq1500005495/982f8770243791581774501290/video_950555_2.m3u8?sign=f25695916ae06c7ea0ee02c47a9aa78f&t=645e46c9'
m3u8="https://video.qianliao.net/43842e87vodtranscq1500005495/2beeac4b243791582056401346/video_950555_2.m3u8?sign=cc33740225b6e7567e7747b68c2a26c9&t=645e6cc1"

m3u8="https://video.qianliao.net/43842e87vodtranscq1500005495/cd65d8203270835008906974310/video_950555_1.m3u8?sign=61d9c9cc6832f902720b856ae8ea2eee&t=645f3a4c"

merged=temp.1.ts
ff="-vf scale=-1:360 -acodec copy -preset veryslow -crf 28"

./xiaoetong -u $m3u8 -ff "$ff" -n $merged -threads 2 -t 20

#ffmpeg -ss 00:47:40 -to 00:52:38  -i "yanzilong.ts" -s 640x360 -crf 28  -acodec copy  yanzilong07.mp4
