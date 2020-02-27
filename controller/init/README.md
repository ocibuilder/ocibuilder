# Init

Init holds the init container for the ocibuilder controller. This container is responsible
for pulling the ocibuilder CRD and storing it in a shared volume to be used by a 
generated ocibuilder build job.

