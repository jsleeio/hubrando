# hubrando

## what?

`hubrando` searches and optionally deletes repository subscriptions in your
Github user account.

## demo

### halp!

```
$ ./hubrando -help
  -delete
    	delete matching repository subscriptions
  -exclude string
    	regular expression with which to exclude subscriptions
  -include string
    	regular expression with which to include subscriptions (default ".")
  -org string
    	constrain operations to repositories within this organization
```

### list all subscriptions to repositories in a single organization

```
$ ./hubrando -org=kubernetes
2018/07/14 14:05:30 found repository subscription kubernetes/kubernetes
2018/07/14 14:05:30 found repository subscription kubernetes/charts
2018/07/14 14:05:30 found repository subscription kubernetes/kube-deploy
```

### list only subscriptions in that org matching a regex

```
$ ./hubrando -org=kubernetes -include=^kube
2018/07/14 14:07:17 found repository subscription kubernetes/kubernetes
2018/07/14 14:07:17 found repository subscription kubernetes/kube-deploy
```

### list only subscriptions in that org NOT matching a regex

```
$ ./hubrando -org=kubernetes -exclude=^charts$
2018/07/14 14:06:43 found repository subscription kubernetes/kubernetes
2018/07/14 14:06:43 found repository subscription kubernetes/kube-deploy
```

### delete subscriptions to all repositories in an org

```
$ ./hubrando -org=kubernetes -delete
2018/07/14 14:08:13 deleted repository subscription kubernetes/kubernetes
2018/07/14 14:08:14 deleted repository subscription kubernetes/charts
2018/07/14 14:08:14 deleted repository subscription kubernetes/kube-deploy
```

