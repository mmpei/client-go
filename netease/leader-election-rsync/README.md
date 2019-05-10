# Leader Election

This is a function that used for stateful set to elect and do some transform if fetch the master lock.

## Running

Run the following three commands in separate terminals. Each terminal needs a unique `id`.

```bash
# first terminal 
go run *.go -kubeconfig=/my/config -logtostderr=true -id=1

# second terminal 
go run *.go -kubeconfig=/my/config -logtostderr=true -id=2

# third terminal
go run *.go -kubeconfig=/my/config -logtostderr=true -id=3
```
> You can ignore the `-kubeconfig` flag if you are running these commands in the Kubernetes cluster.

The one will create a file "/var/leader-election-rsyncfile" if it acquire the lock. and the file will be removed if it lost the lock. So you could do some transform work according to the existing of the file.

Then kill the existing leader. You will see from the terminal outputs that one of the remaining two processes will be elected as the new leader.