(function () {
    function load() {

        const locationInfo = document.getElementById('location');
        const stationInfo = document.getElementById('stationInfo');
        const getAllButton = document.getElementById('getAllButton');
        const getNearestButton = document.getElementById('getNearestButton');
        const clearAll = document.getElementById('clearAll');
        const stationSubUrl = '/customer/station';
        const stationListUrl = location.origin + stationSubUrl;
        const stationNearestUrl = location.origin + stationSubUrl + '/nearest';
        const stationInfoUrl = location.origin + stationSubUrl;
        const stationScooter = location.origin + '/start-trip';


        var map, marker, stations;

        function initMap() {

            map = DG.map('map', {
                center: [48.426, 35.032],
                zoom: 14
            });

            marker = DG.marker([48.426, 35.032], {
                draggable: true
            }).addTo(map).bindLabel('me', {
                static: true
            });

            marker.on('drag', function (e) {
                var lat = e.target._latlng.lat.toFixed(3),
                    lng = e.target._latlng.lng.toFixed(3);

                locationInfo.innerHTML = lat + ', ' + lng;
            });

            stations = new Map();
        }

        async function showStations(sts) {
            stations = new Map();

            sts.forEach((item) => {
                let st = DG.marker([item.latitude, item.longitude], { id: item.id });
                st.on('click', showInfo)
                stations.set(item.id, st);
                st.addTo(map).bindLabel(item.id + "", {
                    static: true
                });
            });
        }

        async function clearStations() {
            stationInfo.innerHTML = "";
            if (stations.size == 0) {
                return
            }

            stations.forEach((val, key) => {
                val.remove();
                stations.delete(key);
            });
        }

        async function showInfo(e) {
            let stid = e.target.options.id;
            let response = await fetch(stationInfoUrl + '/' + stid)

            let data = await response.json();

            stationInfo.innerHTML = "<h4>station info:</h4>" +
                "<p>station id=" + data.id + "</p>" +
                "<p>station name=" + data.name + "</p>" +
                "<p>station is active=" + data.is_active + "</p>" +
                "<a href='/'>show station<" + stationScooter + "/" + data.id + "a>";
        }

        async function getAllStations() {
            let response = await fetch(stationListUrl);
            let data = await response.json();

            await clearStations();
            await showStations(data);
        }

        async function getNearestStation() {
            let latLon = marker.getLatLng();
            let response = await fetch(stationNearestUrl + '?x=' + latLon.lat + '&y=' + latLon.lng);

            let data = await response.json();
            await clearStations();
            await showStations(data);
        }

        DG.then(initMap);
        getAllButton.addEventListener('click', getAllStations);
        getNearestButton.addEventListener('click', getNearestStation);
        clearAll.addEventListener('click', clearStations);
    }

    window.addEventListener('load', function () {
        load();
    });
}())