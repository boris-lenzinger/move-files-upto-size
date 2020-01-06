This small project was written to move files from a partition to another.
I needed to move files but only to free some space. The fastest way to do that
was to code a binary that can move files until it has moved a given amount of
bytes.
This can also apply a filter on the name and require to move the oldest files
first.

To build it:
cd pkg
go build

