export interface SuiteProject {
  name: string;
  website: string;
  repository: string;
}

export const suite: readonly SuiteProject[] = [
  {
    name: 'jira-cli',
    website: 'https://projects.piyushgambhir.com/jira-cli',
    repository: 'https://github.com/piyush-gambhir/jira-cli',
  },
  {
    name: 'cubeapm-cli',
    website: 'https://projects.piyushgambhir.com/cubeapm-cli',
    repository: 'https://github.com/piyush-gambhir/cubeapm-cli',
  },
  {
    name: 'es-cli',
    website: 'https://projects.piyushgambhir.com/es-cli',
    repository: 'https://github.com/piyush-gambhir/es-cli',
  },
  {
    name: 'grafana-cli',
    website: 'https://projects.piyushgambhir.com/grafana-cli',
    repository: 'https://github.com/piyush-gambhir/grafana-cli',
  },
  {
    name: 'jenkins-cli',
    website: 'https://projects.piyushgambhir.com/jenkins-cli',
    repository: 'https://github.com/piyush-gambhir/jenkins-cli',
  },
  {
    name: 'nginxpm-cli',
    website: 'https://projects.piyushgambhir.com/nginxpm-cli',
    repository: 'https://github.com/piyush-gambhir/nginxpm-cli',
  },
  {
    name: 'reckon',
    website: 'https://projects.piyushgambhir.com/reckon',
    repository: 'https://github.com/piyush-gambhir/reckon',
  },
];

export function getOtherSuiteProjects(currentSite: string): SuiteProject[] {
  const currentName = currentSite.split('/').at(-1)?.replace(/\.git$/, '');

  return suite.filter(({ name }) => name !== currentName);
}
