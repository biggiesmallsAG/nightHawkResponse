<p align="center"><img src="https://cloud.githubusercontent.com/assets/17421028/16711172/0e5cc930-4690-11e6-89a4-c0e86d07dd24.png" width="450px"</img></p>

# nightHawk Response

Custom built application for asynchronus forensic data presentation on an Elasticsearch backend.<br>
This application is designed to ingest a Mandiant Redline "collections" file and give flexibility in search/stack and tagging.

**The application was born out of the inability to control multiple investigations (or hundreds of endpoints) in a single pane of glass.**

To ingest redline audits, we created `nightHawk.GO`, a fully fledge GOpher application designed to accompany this framework. The source code to the application is available in this repo, a binary has been compiled and is running inside the iso ready to ingest from first boot.

# Build

**NEW VERSION (1.0.4) is on the way. Start of July, heaps of new features/design !**

01/09/2016: Version 1.0.3 <br>

- User context and user accounts added (login with nighthawk/nighthawk), see wiki article <br>
- Platform statistics and upload information on websockets added <br>
- Delete cases, delete endpoints, delete endpoints from cases added <br>
- Task workflow section on websocket added, see wiki article for guide <br>
- Comments/Tagging are now expandable objects with highlighting enabled <br>
- Comments are now alerts on websocket, see wiki for notes <br>
- Bug in CaseName = Endpoint name resolved <br>
- Zipped audits from Mac/Windows/Linux resolved <br>
- Responsive design upgrade features <br>
- w32system added as audit type <br>

<b>Features:</b>

Video Demonstration: <a href="https://www.youtube.com/watch?v=3bHfAt8bEk8">nightHawk Response Platform</a>

1. Single view endpoint forensics (multiple audit types).<br>
2. Global search.<br>
3. Timelining.<br>
4. Stacking.<br>
5. Tagging.<br>
6. Interactive process tree view.<br>
7. Multiple file upload & Named investigations.<br>

# nightHawk ISO

To make it straight forward for users of nightHawk, we built an ISO with everything setup ready to go. That means you get the following;

1. Latest nightHawk source. <br>
2. CentOS 7 Minimal with core libs needed to operate nightHawk. <br>
3. Nginx and UWSGI setup in reverse proxy (socketed and optimized), SSL enabled. <br>
4. Latest Elasticsearch/Kibana (Kibana is exposed and useable if desired). <br>
5. Sysctrl for all core services. <br>
6. Logging (rotated) for all core services. <br>
7. Configurable system settings, a list of these can be found in the `/opt/nighthawk/etc/nightHawk.json` file. <br>

<b>Starting the system</b>:

Before building your VM with the supplied ISO, take into consideration the following;

1. CPU/RAM. 

_Pending_: Setup the Elastic service to be dual nodes with 1/4 of the allocated system memory per node. This means if you give it 2GB RAM, each ES node will be 512mb and the system will remain with 1GB to operate. <br>

_If you want to set this any different way, ssh into the box and configure your desired way._

2. HDD. 

A minimum of 20GB should be considered. An audit file can be large and therefore its advised you allocate a lot of storage to handle ingesting many collections.

_Pending_: User based storage setup for large scale instances. If you desire to setup extra partitions, you can do this yourself, a few changes can be made to point the ES data storage to your new partition.

<b>Installation</b>:

Download ISO: <a href="https://drive.google.com/open?id=0B-Eozyt1N6W-b1F6RzAtbFJ6d3c">nightHawk v1.0.3</a>

Configure the hardware, mount the ISO into the VM, start the installtion script. <br> 

Once complete, in your browser (Chrome/FireFox), goto; `https://192.168.42.173`. <br>

Log into the system with 'nighthawk/nighthawk' - click "goto site" to get into application <br>

If you need to access Kibana, goto; `https://192.168.42.173:8443`. <br>

If you need to SSH into the box, the login details are; `admin/nightHawk`. <br>

If you want to change the IP address (reflected application wide); `/opt/nighthawk/bin/nighthawkctl set-ip <new_ipaddress>` <br>

Redline Audit Collection Script can be found in the root of this repo. Use this when using the standalone redline collector as this will return the documents you need to populate nightHawk correctly.

<b>Uploading:</b> <br><br>
<b>IMPORTANT</b>:
<b>Creating audit zip file to upload (<u>Redline stand alone collector</u>):</b> <br/>
step_1: Navigate to Sessions\AnalysisSessionX\Audits\<ComputerName>  where X is analysis number which is 1 for most cases. <br/>
step_2: Create zip of folder containing audit files i.e. 20160708085733 <br/>
step_3: Upload 20160708085733.zip <br/><br/>

<b>IMPORTANT</b>: <b>Use exisiting HX audit file (<u>HX collector</u>): </b> FireEye HX audits are an extension ending in .mans. The audit from HX differs from the Redline collector because the .mans that it returns is actually a zip file. This means it can be uploaded directly unlike the Redline audit which you need to follow the instructions above.

Navigate to the "Upload" icon on the nav bar, select an audit .zip (or multiple), a case name (otherwise the system will supply you with one) and submit. If you have used our Redline audit script to build your collection, follow the "Redline Collector" instructions just above. <br>

Once processed, the endpoint will appear in the "Current Investigations" tree node. Under the endpoint you will be presented with all audit types available for that endpoint. The upload feature of this web app spawns pOpen subprocesss that calls the GO application to parse the redline audit and push data into Elasticsearch. There are 2 options for uploading, one is sequential, the other is concurrent. 


_Please Note: Concurrent uploads are limited to 5 at a time and can be resource intensive, if you have an underpowered machine then restrict usage of this feature to 2-3._

<b>Tagging:</b>

You can click on any row in any table (in response view) to tag that data. Once tagged you can view the comments in the comments view.

<b>Elasticsearch:</b>

There are custom mappings (supplied in git root) and advisory comments on the following;

1. <b>Parent/Child relationships:</b>

  Documents are indexed via the GO app as parent/child relation. This was chosen because it is able to give relatively logical path to view documents, ie. the parent is the endpoint name and the children are audit types. Performing aggregations on a parent/child relational document at scale seems to make sense as well. The stacking framework relies on building parents in an array to then get all child document aggregations for certain audit types. <br>

2. <b>Sharding:</b>

  Elasticsearch setups require tuning and proper design recognition. Sharding is important to understand because of the way we are linking parent/child documents. The child is ALWAYS routed to the parent, it cannot exist on its own. This means consideration must be given to how many shards are resident on the index. From what we understand, it may be wise to choose a setup that encorporates many nodes with single shards. To gain performance out of this kind of setup we are working on shard routed searches. 

  We are currently working on designing the best possible configuration for fast searching.
  <br>

3. <b>Scaling:</b>

  This application is designed to scale immensely. From inital design concept, we were able to run it smoothly on a single cpu 2gb ubuntu VM with 3 ES nodes (Macbook Pro), with about 4million+ documents (or 50 endpoints ingested). If going into production, running a setup with 64/128GB RAM and SAS storage, you would be able to maintain a lightning fast response time on document retrival whilst having many analysts working on the application at once.

<b>Considerations:</b>

1. <b>DataTables mixed processing:</b>
   
   There are several audit types ingested that are much to large to return all documents to the table. For example, URL history and Registry may return 15k doc's back to the DOM, rendering this would put strain on the client browser. To combat this, we are using ServerSide processing to page through results of certain audit types. This means you can also search over documents in audit type using Elasticsearch in the backend.

2. <b>Tagging:</b>
   
   Currently we can tag documents and view those comments. We can update them or change them. The analyst is able to give context such as Date/Analyst Name/Comment to the document.

<b>Dependencies (all preinstalled):</b>

`elasticsearch-dsl.py`
`django 1.8 `
`python requests`

<b>To Do:</b>

Process Handles (in progress). <br>
Time selection sliders for time based generators (in progress). <br>
Context menu for Current/Previous investigations. <br>
Tagging context. The tagging system will integrate in a websocket loop for live comments across analyst panes (in progress).<br>
Application context. <br>
Ability to move endpoints between either context. <br>
Potentially redesign node tree to be investigation date driven. <br>
Selective stacking, currently the root node selector is enabled. <br>
Shard routing searches. <br>
Redline Audit script template. <br>
More extensive integration with AngularJS (in progress).<br>
Responsive design. (in progress).<br>
Administrative control page for configuration of core settings (in progress).<br>

<b>Authors & Notes:</b>

We are always looking for likeminded folks who want to contribute to this project, we are by no means web design guru's, if you think we can do something better please request a pull and if we like it, we will merge.<br>

Daniel Eden & 
Roshan Maskey

<b>Credits:</b>

Mandiant Redline devs, AngularJS, Django devs, Angular-DataTables/DataTables, D3 (Bostock), Elasticsearch/ES-dsl.py, jsTree, qTip, GOlang, Python, Fahad Abdulaal (Logo/Video).

# Screenshots:

![alt tag](https://cloud.githubusercontent.com/assets/17421028/16615167/040402da-43b8-11e6-84b9-b25581dbdbfe.png)
![alt tag](https://cloud.githubusercontent.com/assets/17421028/16615162/03a63fb0-43b8-11e6-871a-7515a287043a.png)
![alt tag](https://cloud.githubusercontent.com/assets/17421028/16615166/04000b08-43b8-11e6-8609-442cdba6e5a9.png)
![alt tag](https://cloud.githubusercontent.com/assets/17421028/16615165/03fcebd0-43b8-11e6-8891-17304b796c48.png)
![alt tag](https://cloud.githubusercontent.com/assets/17421028/16615164/03fbac34-43b8-11e6-8a57-33e8b93e94d7.png)
![alt tag](https://cloud.githubusercontent.com/assets/17421028/16615163/03d84870-43b8-11e6-9dce-cbb19dedc118.png)
