# Database

### Creating migrations
To create migrations we are using golang-migrate, it can be installed following the READ.ME in their github repo right
here `https://github.com/golang-migrate/migrate`

After installing `golang-migrate` we can create a new migration by typing
```shell
migrate create -ext {file_extension ex:sql} -dir {our_dir_for_migrations} -seq {sequential id/name for migration}
```

### Database DeadLocks
When performing a queries in transactions we don't have them blocked for the other ones while being used by another
transaction at the same time, and it can cause serious issues when working with concurrent operations.

To lock a single object in a transaction we have to add some arguments to it when performing the query
```postgresql
SELECT *
FROM accounts
WHERE id = $1 
LIMIT 1
FOR NO KEY UPDATE;
```

Above we are saying that we whant to block that row of being updated while another transaction is holding it, and 
`NO KEY UPDATE` that will wait to update the row when another one is holding it, but without fully locking the row.

More info about locks in `https://medium.com/inspiredbrilliance/what-are-database-locks-1aff9117c290`