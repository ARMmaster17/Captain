package org.ARMmaster17.entity;

import io.quarkus.hibernate.orm.panache.PanacheEntity;

import javax.persistence.Entity;

@Entity
public class Plane extends PanacheEntity {
    public String ProviderIdentifier;
    public String Hostname;
    public int Cores;
    public int Memory;
    public String IP;
    public int StorageSize;
}
