- calico报错Calico requires net.ipv4.conf.all.rp_filter to be set to 0 or 1
  
  >  int_dataplane.go 1018: Kernel's RPF check is set to 'loose'. This would allow endpoints to spoof their IP address. Calico requires net.ipv4.conf.all.rp_filter to be set to 0 or 1. If you require loose RPF and you are not concerned about spoofing, this check can be disabled by setting the IgnoreLooseRPF configuration parameter to 'true'.

  kubectl -n kube-system set env daemonset/calico-node FELIX_IGNORELOOSERPF=true 

