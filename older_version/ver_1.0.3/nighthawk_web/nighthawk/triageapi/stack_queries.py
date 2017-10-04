import elasticsearch
from elasticsearch_dsl import Search, Q, A

def GetAuditGenerator(endpoints):
	audit_type = ''
	
	q = []
	a = []

	for k, v in endpoints.iteritems():
		q.append(k)
		for x in v:
			if not v in a:
				a.append(v)

	joined = ' OR '.join([x for x in a[0]])
	generator = q[0]
	
	s = Search()
	s = s[0]

	if generator == 'w32scripting-persistence':
		aggs_gen = A('terms', field='Record.StackPath.raw', size=0)

	elif generator == 'w32prefetch':
		aggs_gen = A('terms', field='Record.ApplicationFileName.raw', size=0)

	elif generator == 'w32network-dns':
		aggs_gen = A('terms', field='Record.RecordName.raw', size=0)

	elif generator in ('urlhistory', 'filedownloadhistory'):
		aggs_gen = A('terms', field='Record.UrlDomain.raw', size=0)

	elif generator == 'w32registryraw':
		aggs_gen = A('terms', field='Record.KeyPath.raw', size=0)

	else:
		aggs_gen = A('terms', field='Record.Name.raw', size=0)

	aggs_endpoint = A('terms', field="ComputerName.raw", size=0)
	s.aggs.bucket('generator', aggs_gen).bucket('endpoint', aggs_endpoint)
	t = Q('query_string', default_field="ComputerName.raw", query=joined)
	query = s.query(t).filter('term', AuditType__Generator=generator)

	return query.to_dict()