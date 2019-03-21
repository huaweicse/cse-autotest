pipeline {
  agent any
  stages {
    stage('dependency') {
      agent any
      steps {
        sh '''echo $WORKSPACE
mkdir -p src/github.com/go-chassis/
cd src/github.com/go-chassis/
git clone https://github.com/go-chassis/go-chassis.git'''
        git(url: 'https://github.com/go-chassis/go-chassis.git', branch: 'master', credentialsId: 'db8a8e53-8f37-4a2f-b104-a637dbad1f45')
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