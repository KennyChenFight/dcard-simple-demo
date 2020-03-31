# Database Setting

For Ubuntu OS Environment

## Step 1: install PostgreSQL

```bash
sudo apt-get install postgresql
```

## Step 2: Creating the database and corresponding user account for the project

Login the postgresql by either one of following command:

1. It will perform OS authentication. It works only if the postgresql is installed in your local machine.

   ```bash
   sudo -u postgres psql -f create_db.sql
   ```

2. The user account should be superuser of postgresql server. Otherwise you are unlikely to have enough privilege.

   ```bash
   psql -h <machine_name> -U <username> -f create_db.sql
   ```

## Step 3: Setup the privilege in the database.

```bash
psql -h <machine_name> -U dcard_admin -f setup_db.sql dcard_db
```

## Step 4: Create the table and foreign key in the database.

```bash
psql -h <machine_name> -U dcard_admin -f create_table.sql dcard_db
psql -h <machine_name> -U dcard_admin dcard_db -f create_fk.sql
```

## Step5: Grant table privilege to the users.

```bash
psql -h <machine_name> -U dcard_admin -f grant_table_privilege.sql dcard_db
```

## Step6(Optional): add testing data
```bash
psql -h <machine_name> -U dcard_admin -f testing_data.sql dcard_db
```
