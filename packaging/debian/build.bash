#!/bin/bash


SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
MM=$SCRIPT_DIR/../..
PHENIX=$MM/phenix


which docker &> /dev/null
DOCKER=$?

if (( $DOCKER == 0 )); then
  echo "Docker detected, so Docker will be used to build minimega and phenix."

  (cd $MM && $MM/docker-build.bash)
  (cd $PHENIX && $PHENIX/docker-build.sh -d)
else
  echo "Docker not detected, so minimega will be built but phenix will not be."

  (cd $MM && $MM/build.bash && $MM/doc.bash)
fi


echo COPYING FILES...

DST=$SCRIPT_DIR/minimega/opt/minimega

mkdir -p $DST

cp -r $MM/bin $DST/
cp -r $MM/doc $DST/
cp -r $MM/lib $DST/

if (( $DOCKER == 0 )); then
  cp -r $PHENIX/bin $DST/
fi

mkdir -p $DST/misc

cp -r $MM/misc/daemon           $DST/misc/
cp -r $MM/misc/uminiccc         $DST/misc/
cp -r $MM/misc/uminirouter      $DST/misc/
cp -r $MM/misc/vmbetter_configs $DST/misc/
cp -r $MM/misc/web              $DST/misc/

DOCS=$SCRIPT_DIR/minimega/usr/share/doc/minimega

mkdir -p $DOCS

cp $MM/LICENSE     $DOCS/
cp -r $MM/LICENSES $DOCS/

echo COPIED FILES


echo BUILDING DEB PACKAGE...

(cd $SCRIPT_DIR && fakeroot dpkg-deb -b minimega)

echo DONE BUILDING DEB PACKAGE
