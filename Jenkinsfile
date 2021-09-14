#!groovy

milestone 0
timestamps {
    node('docker') {
        def commitHash = checkout(scm).GIT_COMMIT

        docker.withRegistry('https://harbor.cyverse.org', 'jenkins-harbor-credentials') {
            def dockerImage
            stage('Build') {
                milestone 50
                dockerImage = docker.build("harbor.cyverse.org/de/de-job-killer:${env.BUILD_TAG}", "--build-arg git_commit=${commitHash} .")
                milestone 51
                dockerImage.push()
            }
            stage('Docker Push') {
                milestone 100
                dockerImage.push("${env.BRANCH_NAME}")
                // Retag to 'qa' if this is master/main (keep both so when it switches this keeps working)
                if ( "${env.BRANCH_NAME}" == "master" || "${env.BRANCH_NAME}" == "main" ) {
                    dockerImage.push("qa")
                }
                milestone 101
            }
        }
    }
}
