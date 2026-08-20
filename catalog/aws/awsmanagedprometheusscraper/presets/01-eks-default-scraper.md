# EKS Cluster Scraper (AWS Default Configuration)

The day-1 shape: point the scraper at your EKS cluster and AMP workspace and leave `scrapeConfiguration` unset — the modules resolve AWS's published default (kubelet, cAdvisor, pod service discovery) at deploy. Metrics flow with zero agents installed.
