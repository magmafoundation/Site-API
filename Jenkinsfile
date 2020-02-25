pipeline {
  agent any
  stages {
    stage('Docker-Build') {
      steps {
           script {
            dockerImage = docker.build("hexeption/magma-api:$BUILD_NUMBER", "--no-cache . ")
            dockerImage.push()
          }
        }
      }
    }
    stage('Deploy') {
      steps {
        sh "kubectl set image deployment magma-api magma-api=hexeption/magma-api:$BUILD_NUMBER --namespace magma-api"
      }
    }
 }
