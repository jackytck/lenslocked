#!/bin/bash
sudo yum install postgresql-server postgresql-contrib
sudo postgresql-setup initdb
sudo systemctl start postgresql
sudo systemctl enable postgresql
sudo -u postgres psql
ALTER USER postgres WITH ENCRYPTED PASSWORD 'natnat';
CREATE DATABASE lenslocked_prod;
sudo su - postgres
cd data
vi pg_hba.conf
logout
sudo systemctl restart postgresql