package org.ARMmaster17.service.provider.proxmox;

import it.corsinvest.proxmoxve.api.PveClient;
import it.corsinvest.proxmoxve.api.Result;
import org.ARMmaster17.service.provider.ProviderBackendService;
import org.eclipse.microprofile.config.inject.ConfigProperty;

import javax.enterprise.context.ApplicationScoped;
import javax.inject.Inject;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.util.HashMap;
import java.util.Map;

//@LookupIfProperty(name = "provider.proxmox.enabled", stringValue = "true")
@ApplicationScoped
public class ProxmoxProvider implements ProviderBackendService {

    static int TASK_RECHECK_INTERVAL_MS = 100;
    static int TASK_RECHECK_TIMEOUT_MS = 1000 * 20;

    @Inject
    PveClient pveClient;

    @ConfigProperty(name = "provider.proxmox.default.ostemplate")
    String osTemplate;

    @ConfigProperty(name = "provider.all.dns")
    String nameserver;

    @ConfigProperty(name = "provider.proxmox.default.bridge")
    String bridge;

    @ConfigProperty(name = "provider.proxmox.default.node")
    String node;

    @ConfigProperty(name = "provider.all.gateway")
    String gateway;

    @ConfigProperty(name = "provider.all.mtu")
    int mtu;

    @ConfigProperty(name = "provider.proxmox.default.ostype")
    String ostype;

    @ConfigProperty(name = "provider.proxmox.default.storageshare")
    String storageShare;

    @ConfigProperty(name = "provider.all.searchdomain")
    String searchdomain;

    @Override
    public String CreateMachine(String hostname, int cores, int memory, String ip, int storageSize) throws Exception {
        int vmid = getNextVMID();
        String arch = "amd64";
        Float bwlimit = null;
        String cmode = "tty";
        Boolean console = true;
        Float cpulimit = null;
        Integer cpuunits = null;
        Boolean debug = null;
        String description = "Managed by the Captain stack.";
        String features = null;
        Boolean force = null;
        String hookscript = null;
        Boolean ignore_unpack_errors = null;
        String lock_ = null;
        Map<Integer, String> mpN = null;
        Map<Integer, String> netN = new HashMap<>();
        netN.put(0, URLEncoder.encode(String.format("name=eth0,bridge=%s,firewall=0,gw=%s,ip=%s,mtu=%d", bridge, gateway, ip, mtu), StandardCharsets.UTF_8));
        Boolean onboot = true;
        String ostemplate = URLEncoder.encode(osTemplate, StandardCharsets.UTF_8);
        String password = "captain"; // TODO: Not make this hard-coded?
        String pool = null;
        Boolean protection = false;
        Boolean restore = null;
        String rootfs = String.format("%d", storageSize);
        String ssh_public_keys = null; // TODO: Pull from system configuration.
        Boolean start = true;
        String startup = null;
        String storage = storageShare;
        int swap = memory / 2;
        String tags = null;
        Boolean template = false;
        String timezone = "host";
        int tty = 2;
        Boolean unique = null;
        Boolean unprivileged = true;
        Map<Integer, String> unusedN = null;

        pveClient.setDebugLevel(99);

        Result result = pveClient.getNodes().get(node).getLxc().createVm(ostemplate, vmid, arch, bwlimit, cmode,
                console, cores, cpulimit, cpuunits, debug, description, features, force, hookscript, hostname,
                ignore_unpack_errors, lock_, memory, mpN, nameserver, netN, onboot, ostype, password, pool, protection,
                restore, rootfs, searchdomain, ssh_public_keys, start, startup, storage, swap, tags, template, timezone,
                tty, unique, unprivileged, unusedN);

        if(!result.isSuccessStatusCode()) {
            throw new Exception(String.format("%s: %s", result.getError(), result.getReasonPhrase()));
        }
        String taskIdentifier = (String)result.getResponse().get("data");
        pveClient.waitForTaskToFinish(taskIdentifier, TASK_RECHECK_INTERVAL_MS, TASK_RECHECK_TIMEOUT_MS);
        String taskResult = pveClient.getExitStatusTask(taskIdentifier);
        if(taskResult != "OK") {
            throw new Exception("Failed to create new VM: " + taskResult);
        }

        return getVMIDFromStatus(taskIdentifier);
    }

    @Override
    public void DestroyMachine(String providerID) throws Exception {
        stopLXCAndWait(providerID);
        Result result = pveClient.getNodes().get(node).getLxc().get(providerID).destroyVm(true, true, true);
        if(!result.isSuccessStatusCode()) {
            throw new Exception(String.format("%s: %s", result.getError(), result.getReasonPhrase()));
        }
    }

    @Override
    public void Ping() throws Exception {
//        if (proxmoxProvider.cluster().nextId() == -1) {
//            throw new Exception("Unable to connect to Proxmox cluster");
//        }
    }

    int getNextVMID() {
        return pveClient.getCluster().getNextid().nextid().getResponse().getInt("data");
    }

    String getVMIDFromStatus(String result) {
        return result.split(":")[6];
    }

    void stopLXCAndWait(String providerID) throws Exception {
        Result result = pveClient.getNodes().get(node).getLxc().get(providerID).getStatus().getStop().vmStop();
        if(!result.isSuccessStatusCode()) {
            throw new Exception(String.format("%s: %s", result.getError(), result.getReasonPhrase()));
        }
        pveClient.waitForTaskToFinish(result.getResponse().getString("data"), TASK_RECHECK_INTERVAL_MS, TASK_RECHECK_TIMEOUT_MS);
    }
}
