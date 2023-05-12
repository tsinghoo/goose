# build
# go build
set -x
#download
m3u8='https://video.qianliao.net/43842e87vodtranscq1500005495/982f8770243791581774501290/video_950555_2.m3u8?sign=f25695916ae06c7ea0ee02c47a9aa78f&t=645e46c9'

merged=temp.1.ts
ff="-vf scale=-1:240 -acodec copy -preset veryslow -crf 28"

./xiaoetong -u $m3u8 -ff "$ff" -n $merged -threads 2 -t 11