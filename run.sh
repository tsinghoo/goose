# build
# go build
set -x
#download
#hu
m3u8='https://video.qianliao.net/43842e87vodtranscq1500005495/982f8770243791581774501290/video_950555_2.m3u8?sign=f25695916ae06c7ea0ee02c47a9aa78f&t=645e46c9'
#Yuan
m3u8="https://video.qianliao.net/43842e87vodtranscq1500005495/2beeac4b243791582056401346/video_950555_2.m3u8?sign=cc33740225b6e7567e7747b68c2a26c9&t=645e6cc1"
#shen
m3u8="https://video.qianliao.net/43842e87vodtranscq1500005495/c510d6753270835008954474999/video_950555_1.m3u8?sign=091e959f2af899ec90e1bf36bc5c2786&t=645f34af"
#xiong
m3u8="https://video.qianliao.net/43842e87vodtranscq1500005495/cd65d8203270835008906974310/video_950555_1.m3u8?sign=61d9c9cc6832f902720b856ae8ea2eee&t=645f3a4c"
m3u8="https://mps-trans.yzcdn.cn/multi_trans_hls_hd/67j57BpSVwY-VWZDvLklgv96sG8bcS8_3Lymog_HD.m3u8?sign=257ace0911d3831f37a5db688b610034&t=6466d928"
m3u8="https://video.qianliao.net/43842e87vodtranscq1500005495/db125433243791581751601387/video_950555_1.m3u8?sign=99dea2a2fad03d5e6900b1c9539d154c&t=6466440a"
m3u8="https://hss4.dnvodcdn.me/ppot/_definst_/mp4:s4/vod/jq-tmz-yueyu-026A421FC.mp4/chunklist.m3u8?dnvodendtime=1685194521&dnvodhash=Fpf1dj7y3T8Zq3mIJoxYB20Wgw4uGYd5R1ZKISdCPfI=&dnvodCustomParameter=0_139.180.146.177.SG_1&lb=1efb7e8350aaa9dbce7fbbe69556395a&us=1&vv=e7a2e22c7798ae145ad075804812454b&pub=CJOuDJ0nE34mDIuvCJHVKqTVCJCvBZ4uC2unD3OkCJStNpGqPJOuEMGsPZbcPZGnDJ5YOMOnDcOmOZ1XCJ9XEJCvNpPcCp1YOZCsOcGuCJTYDpLcPJCtP6KuDMCvDJLZDJba"
m3u8="https://encrypt-k-vod.xet.tech/522ff1e0vodcq1252524126/6411fc3b3270835010198160916/playlist_eof.m3u8?sign=70abd334f3f32f902a644fb20cc2d07b&t=64a41c00&us=qXybDdrmAO&time=1688433473961&uuid=u_6483d088220d6_dr03V4rG5f"
m3u8="https://video.qianliao.net/43842e87vodtranscq1500005495/0facd5a2387702294794384097/adp.950555.m3u8?t=65299b0c&sign=3a6ab0769b7a44ddf9f6d6724eb0a051"
m3u8="https://video.qianliao.net/43842e87vodtranscq1500005495/0facd5a2387702294794384097/video_950555_0.m3u8?sign=3a6ab0769b7a44ddf9f6d6724eb0a051&t=65299b0c"
m3u8="https://video.qianliao.net/43842e87vodtranscq1500005495/b70962d4387702299113206595/video_950555_0.m3u8?sign=b6aba39fe4905a845c90784b2a0d9458&t=6529a557"
m3u8="https://video.qianliao.net/43842e87vodtranscq1500005495/e520d3e9387702297491947759/video_950555_0.m3u8?sign=7f88868c16ceac708de846f6111d828f&t=6529a7dc"
m3u8="https://video2.qianliao.net/43a5a542vodtranscq1500020547/d76f400b5576678020618140668/video_1440284_0.m3u8?sign=363f8d0f3a42cbbacc00c8447d35164b&t=652a3ec8"
#video_950555_0.m3u8?sign=3a6ab0769b7a44ddf9f6d6724eb0a051&t=65299b0c

ff="-vf scale=-1:480 -to 00:00:10 -c:v libx264 -acodec copy -preset veryslow -crf 28"
ff="-vf scale=-1:480 -c:v libx264 -acodec copy -preset veryslow -crf 28"


now=`date "+%Y-%m-%d %H:%M:%S"`
echo "$now start" >> run.log
merged=qixue4

./goose -u $m3u8 -ff "$ff" -n $merged -threads 2 -t 0 >> run.log


now=`date "+%Y-%m-%d %H:%M:%S"`
echo "$now end" >> run.log
#ffmpeg -ss 00:47:40 -to 00:52:38  -i "yanzilong.ts" -s 640x360 -crf 28  -acodec copy  yanzilong07.mp4
#ffmpeg -i huxing.ts -vf scale=-1:360 -acodec copy -preset veryslow -crf 28 huxing.mp4
