async function getUserWords() {
    var myHeaders = new Headers();
    myHeaders.append("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDk4OTk0MTEsImlhdCI6MTcwMTI1OTQxMSwiaXAiOiI6OjEiLCJ1aWQiOiJhMDc4ZmQ0Mi0zNWVlLTQ5N2MtYmRiOC1mNjlkOWM0NjUxMzEifQ.Qby8DQyyZuAovCig_DyzkZvVRG0bBaSlxD_eFytz224");

    var requestOptions = {
        method: 'GET',
        headers: myHeaders,
        redirect: 'follow'
    };

    try {
        const response = await fetch("http://localhost:8080/user/words", requestOptions);
        const result = await response.json(); // Parse the response as JSON
        return result;
    } catch (error) {
        console.error('Error:', error);
        return null; // or handle the error accordingly
    }
}

async function sendWordRequest(data) {
    const body = JSON.stringify(data);

    try {
        const response = await fetch('http://localhost:8080/words/add', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDk4OTk1MjUsImlhdCI6MTcwMTI1OTUyNSwiaXAiOiI6OjEiLCJ1aWQiOiJhMDc4ZmQ0Mi0zNWVlLTQ5N2MtYmRiOC1mNjlkOWM0NjUxMzEifQ.rbjGW7jGv61bviTINTGvzum4FpGfnhdSNRqNPtTwRIw',
            },
            body: body,
        });

        const responseData = await response.json();
        console.log(responseData);
    } catch (error) {
        console.error('Error:', error);
    }
}

async function getDefinition(word) {
    const apiUrl = `https://api.datamuse.com/words?sp=${word}&md=d`;

    try {
        const response = await fetch(apiUrl);

        if (!response.ok) {
            throw new Error('Failed to fetch definition.');
        }

        const data = await response.json();
        const definition = data.length > 0
            ? data[0].defs && data[0].defs.length > 0
                ? data[0].defs[0]
                : 'Definition not found.'
            : 'Definition not found.';

        return definition
    } catch (error) {
        console.error('Error fetching definition:', error);
    }
}

function convertJsonDataToObject(jsonData) {
    const result = {};

    jsonData.words.forEach((item) => {
        const wordKey = item.word.word;
        const langKey = item.word.lang;

        if (!result[wordKey]) {
            result[wordKey] = {};
        }

        result[wordKey][langKey] = {
            example: item.word.example,
            image_url: item.word.image_url,
            link: item.word.link
        };

        if (!result[wordKey].definitions) {
            result[wordKey].definitions = [];
        }

        item.definitions.forEach((definition) => {
            result[wordKey].definitions.push({
                lang: definition.lang,
                definition: definition.definition
            });
        });
    });

    return result;
}