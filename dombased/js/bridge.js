function __dombasedxssfinder_vulns_push(sources, sinkLabel) {
    let results = [];
    for (const source of sources) {
        let r = {
            url: location.href,
            source,
            sink: __dombasedxssfinder_get_sink(sinkLabel)
        }
        results.push(r);
        console.debug('result', r);
    }
    // Runtime.addBinding ==> PushDomVul
    window.PushDomVul(JSON.stringify(results));
}