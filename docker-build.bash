#!/bin/bash -e


which docker &> /dev/null

if (( $? )); then
  echo "Docker must be installed (and in your PATH) to use this build script. Exiting."
  exit 1
fi


USER_UID=$(id -u)
USERNAME=builder


docker build -t minimega:builder -f - . <<EOF
FROM golang:1.12

RUN groupadd --gid $USER_UID $USERNAME \
  && useradd -s /bin/bash --uid $USER_UID --gid $USER_UID -m $USERNAME

RUN apt update && apt install -y libpcap-dev
EOF


echo BUILDING MINIMEGA...

docker run -it --rm -v $(pwd):/minimega -w /minimega -u $USERNAME minimega:builder ./clean.bash
docker run -it --rm -v $(pwd):/minimega -w /minimega -u $USERNAME minimega:builder ./build.bash
docker run -it --rm -v $(pwd):/minimega -w /minimega -u $USERNAME minimega:builder ./doc.bash

bin/pyapigen bin/minimega

echo DONE BUILDING MINIMEGA


(cd phenix && ./docker-build.sh -d)
cp phenix/bin/phenix bin/phenix