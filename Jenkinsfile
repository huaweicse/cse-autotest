pipeline {
  agent any
  stages {
    stage('build') {
      agent any
      steps {
        sh '''rm -rf $WORKSPACE/src/github.com/go-chassis/
echo $WORKSPACE
mkdir -p $WORKSPACE/src/github.com/go-chassis/
cd $WORKSPACE/src/github.com/go-chassis/
git clone https://github.com/go-chassis/go-chassis.git



export GOPATH=$WORKSPACE
export PATH=/usr/local/go/bin:$PATH
cd $WORKSPACE/src/github.com/go-chassis/go-chassis
GO111MODULE=on go mod vendor
mv $WORKSPACE/src/github.com/go-chassis/go-chassis/vendor $WORKSPACE/src/
cd $WORKSPACE/scripts
go get github.com/huaweicse/auth
go get github.com/huaweicse/cse-collector
bash build_gosdk_demo_image.sh'''
      }
    }
  }
  environment {
    SDKAT_SWR_ADDR = 'swr.cn-east-2.myhuaweicloud.com'
    SDKAT_SWR_ORG = 'tian'
  }
}