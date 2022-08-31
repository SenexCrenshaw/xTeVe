class TVLogosFile {
    /**
     * @param currentLogoUrl Current channel logo to set initial value to.
     * @returns Array of, sequentially:
     * 1) Container of the picker.
     * 2) Input field to type at and get choice from.
     * 3) Datalist containing every option.
     */
    newTvLogosIdPicker(currentLogoUrl: string): [HTMLDivElement, HTMLInputElement, HTMLDataListElement] {
        const container = document.createElement('div');
        const input = document.createElement('input');
        input.setAttribute('type', 'text');

        // Initially, set value to blank if input is empty
        input.value = (currentLogoUrl) ? currentLogoUrl : '';

        // When input lose focus or take a value, if it's empty, set value to blank
        input.addEventListener('blur', setFallbackValue);
        input.addEventListener('change', setFallbackValue);
        function setFallbackValue(evt: Event) {
            const target = evt.target as HTMLInputElement;
            target.value = (target.value) ? target.value : (SERVER['settings']['tlsMode']) ? 'https://' + SERVER['clientInfo']['DVR'] + '/web/img/tv-test-pattern.png' : 'http://' + SERVER['clientInfo']['DVR'] + '/web/img/tv-test-pattern.png';
        }

        container.appendChild(input);

        const datalist = document.createElement('datalist');

        const option = document.createElement('option');
        option.setAttribute('value', (SERVER['settings']['tlsMode']) ? 'https://' + SERVER['clientInfo']['DVR'] + '/web/img/tv-test-pattern.png' : 'http://' + SERVER['clientInfo']['DVR'] + '/web/img/tv-test-pattern.png');
        option.innerText = 'Default non-custom logo';
        datalist.appendChild(option);


        if (SERVER['tvlogos']) {
            SERVER['tvlogos']["LogoInformation"].forEach((logo) => {
                if (SERVER["settings"]["logosCountry"] === "All" ||
                    logo["path"].startsWith("misc") ||
                    logo["country"] === SERVER["settings"]["logosCountry"]
                ) {
                    const option = document.createElement('option');
                    option.setAttribute('value', logo["path"]);
                    option.innerText = logo["filename"];
                    datalist.appendChild(option);
                }
            });

        }

        container.appendChild(datalist);

        return [container, input, datalist];
    }

    newTvLogosCountryPicker(currentCountry: string): [HTMLDivElement, HTMLInputElement, HTMLDataListElement] {
        const container = document.createElement('div');

        const input = document.createElement('input');
        input.value = (currentCountry) ? currentCountry : '';
        container.appendChild(input);

        const datalist = document.createElement('datalist');

        const option = document.createElement('option');
        option.setAttribute('value', 'All');
        datalist.appendChild(option);

        if (SERVER['tvlogos']) {
            let countries: (string | undefined)[] = [];
            SERVER['tvlogos']["LogoInformation"].forEach((logo) => {
                if (!logo["path"].startsWith("misc") &&
                    countries.indexOf(logo["country"]) === -1) {
                    countries.push(logo["country"]);
                }
            });
            countries.sort();

            countries.forEach((country) => {
                const option = document.createElement('option');
                option.setAttribute('value', country);
                datalist.appendChild(option);
            });
        }
        container.appendChild(datalist);

        return [container, input, datalist];
    }

    return: any;
}
