pipeline {
  agent any
  stages {
    stage('dependency') {
      agent any
      steps {
        sh '''echo $WORKSPACE
mkdir -p src/github.com/go-chassis/
cd src/github.com/go-chassis/
'''
        git(url: 'git@github.com:go-chassis/go-chassis.git', branch: 'master', credentialsId: 'db8a8e53-8f37-4a2f-b104-a637dbad1f45', poll: true)
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