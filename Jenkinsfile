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
    stage('push image') {
      steps {
        sh '''export login_key=`printf "$AK" | openssl dgst -binary -sha256 -hmac "$SK" | od -An -vtx1 | sed \'s/[ \\n]//g\' | sed \'N;s/\\n//\'`


export SDKAT_SWR_LOGIN_CMD="docker login -u $REGION@$AK -p $login_key $SDKAT_SWR_ADDR"



cd scripts

bash push_gosdk_demo_to_huaweicloud.sh'''
      }
    }
    stage('upgrade AOS') {
      steps {
        sh '''cd scripts 
bash upgrade_stack gosdk'''
      }
    }
  }
  parameters {
    string(defaultValue: 'swr.cn-east-2.myhuaweicloud.com', description: '', name: 'SDKAT_SWR_ADDR')
    string(defaultValue: 'go-chassis', description: '', name: 'SDKAT_SWR_ORG')
    string(defaultValue: 'cn-north-1', description: '', name: 'REGION')
    string(defaultValue: '', description: '', name: 'AK')
    string(defaultValue: '', description: '', name: 'SK')
  }
}