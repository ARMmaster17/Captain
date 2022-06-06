package org.ARMmaster17.service.provider.proxmox;

import it.corsinvest.proxmoxve.api.PveClient;
import org.eclipse.microprofile.config.inject.ConfigProperty;

import javax.enterprise.context.ApplicationScoped;
import javax.enterprise.context.Dependent;
import javax.enterprise.inject.Default;
import javax.inject.Singleton;
import javax.ws.rs.Produces;

@ApplicationScoped
public class ProxmoxVEClientConfig {

    @ConfigProperty(name = "provider.proxmox.endpoint")
    String endpoint;

    @ConfigProperty(name = "provider.proxmox.realm")
    String realm;

    @ConfigProperty(name = "provider.proxmox.username")
    String username;

    @ConfigProperty(name = "provider.proxmox.password")
    String password;

    @Produces
    @Default
    PveClient pveClient() {
        try {
            PveClient client = new PveClient(endpoint, 8006);
            if (!client.login(username, password, realm)) {
                throw new Exception("Unable to connect to Proxmox cluster");
            }
            return client;
        } catch (Exception e) {
            return null;
        }

    }
}
