pipeline {
  agent any
  stages {
    stage('checkout code') {
      steps {
        git(url: 'https://github.com/almazik77/tg-bot', branch: 'dev', changelog: true)
      }
    }

  }
}