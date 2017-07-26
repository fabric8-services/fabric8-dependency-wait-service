#!/usr/bin/groovy
@Library('github.com/rawlingsj/fabric8-pipeline-library@master')
def dummy
goNode{
  dockerNode{
    if (env.BRANCH_NAME.startsWith('PR-')) {
      goCI{
        githubOrganisation = 'fabric8-services'
        dockerOrganisation = 'fabric8'
        project = 'fabric8-dependency-wait-service'
        makeTarget = 'clean test cross'
      }
    } else if (env.BRANCH_NAME.equals('master')) {
      def v = goRelease{
        githubOrganisation = 'fabric8-services'
        dockerOrganisation = 'fabric8'
        project = 'fabric8-dependency-wait-service'
      }

      stage ('Update downstream dependencies') {
        updateDownstreamDependencies(v)
      }
    }
  }
}

def updateDownstreamDependencies(v) {
  pushPomPropertyChangePR {
    propertyName = 'dependency-wait-service.version'
    projects = [
            'fabric8io/fabric8-platform',
    ]
    version = v
  }
}
