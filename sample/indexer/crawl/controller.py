import http.client
import json
import sys
import getopt
import queue
import threading

def print_usage():
    print('conroller.py <options>\n')
    print('options: \n')
    print('-h                    print this help and quit\n')
    print('-c --crawlerhost=URL  url to reach the crawler service\n')
    print('-s --seedurl=URL      root url to start the crawling\n')
    print('-r --repodir=PATH     full path to the directory where crawled documents are stored\n')

crawler_url = ''
crawl_seed_url = ''
crawl_repository = '.'

try:
    opts, args = getopt.getopt(argv,"hc:s:r",["crawlerhost=","seedurl=", "repodir="])
except getopt.GetoptError:
    print_usage()
    sys.exit(2)

for opt, arg in opts:
    if opt == '-h':
        print_usage()
        sys.exit()
    elif opt in ("-c", "--crawlerhost"):
        crawler_url = arg
    elif opt in ("-s", "--seedurl"):
        crawl_seed_url = arg
    elif opt in ("-r", "--repodir"):
        crawl_repository = arg

crawler_host = crawler_url[crawler_url.find('://')+3:]

def urlfile(url):
    return url[url.rfind("/")+1:]

def getconnection():
    if crawler_url.startswith('https'):
        return http.client.HTTPSConnection(crawler_host)
    else:
        return http.client.HTTPConnection(crawler_host)

headers = {'Content-type': 'application/json'}
urls_q = queue.Queue()
docs_q = queue.Queue()

def crawl(urls):
    resources = {}
    for url in urls:
        k = urlfile(url)
        resources[k] = url
        docs_q.put([k, url])

    request = {"repository": "file",
               "repositoryConfig": crawl_repository,
               "resources": {},
               "context": {}}
    request["resources"] = resources
    json_req = json.dumps(request)
    connection = None
    try:
        connection = getconnection()
        connection.request('POST', '/crawl', json_req, headers)
        response = connection.getresponse()
    except Exception as ex:
        print(ex)
    finally:
        if connection:
            connection.close()

def crawl_job():
    while True:
        try:
            urls = urls_q.get(False)
        except queue.Error:
            continue

t = threading.Thread(target=get_url, args = (q,u))
t.daemon = True
t.start()
