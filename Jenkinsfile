pipeline {
  agent any
  stages {
    stage('dependency') {
      agent any
      steps {
        sh '''echo $WORKSPACE
mkdir -p src/github.com/go-chassis/
cd src/github.com/go-chassis/
git clone https://github.com/go-chassis/go-chassis.git
export GOPATH=$WORKSPACE
GO111MODULE=on go mod download'''
        sh '''export GOPATH=$WORKSPACE
GO111MODULE=on go mod download'''
      }
    }
    stage('build') {
      steps {
        sh 'bash build_all.sh'
      }
    }
  }
}