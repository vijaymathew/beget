How to run the "sample" distributed crawler infrastructure.
===========================================================

Start Beget in http mode:

$ ./beget -httpd true

Start the simpleHTTP repository (i.e run httprepo.py in this directory, needs flask):

$ python httprepo.py

Start the crawl controller (needs python3):

$ python3 controller.py -c http://localhost:8080 -u urls.txt -t simpleHTTP -r http://localhost:5050/doc

Now you should see Beget crawling the resources listed in urls.txt and pushing the data to the `httprepo` service.