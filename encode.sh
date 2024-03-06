# build
# go build
cd=`pwd`
echo cd=$cd
cd $cd
#download
ff="-vf scale=-1:480 -to 00:00:10 -c:v libx264 -acodec copy -preset veryslow -crf 28"
ff="-vf scale=-1:480 -c:v libx264 -acodec copy -preset veryslow -crf 28"
ff="-c:v libx264 -acodec copy -preset veryslow -crf 28"

info(){
 now=`date +%Y/%m/%d_%H-%M-%S`
 echo [$now] $@
}
encode(){
  info encode start
  #tss=`ls *.ts`
  tss=`ls ~/flv/baidu/*10.mp4`
  for ts in $tss;do
   if [ -f "$ts.mp4" ]; then
     info $ts skipped
   else
     set -x
     ffmpeg -i $ts $ff $ts.mp4
     set +x
   fi
  done
  info encode end
}

run(){
  encode >> log.encode 2>&1
  sleep 1
}

run

