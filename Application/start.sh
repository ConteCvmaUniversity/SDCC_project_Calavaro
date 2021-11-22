if [ $# -eq 0 ]
  then
    docker-compose up
  else
   case "$1" in
    -b)
        docker-compose up --build
   esac
   
fi

