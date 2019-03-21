pipeline {
  agent any
  stages {
    stage('dependency') {
      agent any
      steps {
        sh '''rm -rf src/github.com/go-chassis/
echo $WORKSPACE
mkdir -p src/github.com/go-chassis/
cd src/github.com/go-chassis/
git clone https://github.com/go-chassis/go-chassis.git
export GOPATH=$WORKSPACE
export PATH=/usr/local/go/bin:$PATH
GO111MODULE=on go mod download'''
      }
    }
    stage('build') {
      steps {
        sh '''cd $WORKSPACE/scripts
bash build_gosdk_demo_image.sh'''
      }
    }
  }
}