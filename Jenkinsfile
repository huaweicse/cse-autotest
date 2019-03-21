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
        git(url: 'git@github.com:go-chassis/go-chassis.git', branch: 'master', credentialsId: 'bedce822-c43b-4538-9f19-a42a6764b53a')
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