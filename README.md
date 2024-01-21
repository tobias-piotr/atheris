# 🐉 Atheris

Taking control of the APIs.

The basic idea of this project was to provide a single proxy/gateway that would
listen to all requests at a single port and then route them to the appropriate
service. In the meantime, it would also save the requests and responses, so
that they can be analyzed and compared later on. Proxy server was developed and
there is also a CLI for viewing, renaming and deleting requests, but I got
bored somewhere in the middle, so it will probably wait for more features to be
added.
