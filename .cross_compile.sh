set -e
###
 # @Author: Vincent Young
 # @Date: 2023-10-30 23:04:24
 # @LastEditors: Vincent Young
 # @LastEditTime: 2023-10-30 23:04:35
 # @FilePath: /AppStorePrice/.cross_compile.sh
 # @Telegram: https://t.me/missuo
 # 
 # Copyright © 2023 by Vincent, All Rights Reserved. 
### 
DIST_PREFIX="as_price"
DEBUG_MODE=${2}
TARGET_DIR="dist"
PLATFORMS="darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/386 windows/amd64"

rm -rf ${TARGET_DIR}
mkdir ${TARGET_DIR}

for pl in ${PLATFORMS}; do
    export GOOS=$(echo ${pl} | cut -d'/' -f1)
    export GOARCH=$(echo ${pl} | cut -d'/' -f2)
    export TARGET=${TARGET_DIR}/${DIST_PREFIX}_${GOOS}_${GOARCH}
    if [ "${GOOS}" == "windows" ]; then
        export TARGET=${TARGET_DIR}/${DIST_PREFIX}_${GOOS}_${GOARCH}.exe
    fi

    echo "build => ${TARGET}"
    if [ "${DEBUG_MODE}" == "debug" ]; then
        CGO_ENABLED=0 go build -trimpath -gcflags "all=-N -l" -o ${TARGET} \
            -ldflags "-w -s" main.go
    else
        CGO_ENABLED=0 go build -trimpath -o ${TARGET} \
            -ldflags "-w -s" main.go
    fi
done