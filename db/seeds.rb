# This file should contain all the record creation needed to seed the database with its default values.
# The data can then be loaded with the bin/rails db:seed command (or created alongside the database with db:setup).
#
# Examples:
#
#   movies = Movie.create([{ name: 'Star Wars' }, { name: 'Lord of the Rings' }])
#   Character.create(name: 'Luke', movie: movies.first)
Machine.create(hostname: "test1.example.com", ip_address: "10.1.0.1", vmid: 103, cpu: 1, ram: 256, disk: 8)
Machine.create(hostname: "test2.example.com", ip_address: "10.1.0.2", vmid: 104, cpu: 1, ram: 256, disk: 8)
Machine.create(hostname: "test3.example.com", ip_address: "10.1.0.3", vmid: 105, cpu: 1, ram: 256, disk: 8)
Machine.create(hostname: "test4.example.com", ip_address: "10.1.0.4", vmid: 106, cpu: 1, ram: 256, disk: 8)