pipeline {
  agent any
  stages {
      stage('Docker-Build') {
          steps {
              script {
                  dockerImage = docker.build("docker.hexeption.dev/api/magma-api:$BUILD_NUMBER", "--no-cache . ")
                  dockerImage.push()
              }
          }
      }
  }
 }
