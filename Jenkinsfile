pipeline {
  agent any
  stages {
    stage('build') {
      agent any
      steps {
        git(url: 'git@github.com:go-chassis/go-chassis.git', branch: 'master')
      }
    }
  }
}