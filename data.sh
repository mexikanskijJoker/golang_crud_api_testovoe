# url="http://localhost:8080/api/v1/create"
# headers="Content-Type: application/json"
# data_list=(
#     '{"Group": "Dua Lipa", "Song": "Houdini"}'
#     '{"Group": "Coldplay", "Song": "Yellow"}'
#     '{"Group": "The Beatles", "Song": "Hey Jude"}'
#     '{"Group": "Queen", "Song": "Bohemian Rhapsody"}'
#     '{"Group": "Nirvana", "Song": "Smells Like Teen Spirit"}'
#     '{"Group": "Radiohead", "Song": "Creep"}'
#     '{"Group": "Adele", "Song": "Hello"}'
#     '{"Group": "The Rolling Stones", "Song": "Paint It Black"}'
#     '{"Group": "Pink Floyd", "Song": "Comfortably Numb"}'
#     '{"Group": "Led Zeppelin", "Song": "Stairway to Heaven"}'
#     '{"Group": "U2", "Song": "With or Without You"}'
#     '{"Group": "The Police", "Song": "Every Breath You Take"}'
#     '{"Group": "The Clash", "Song": "London Calling"}'
#     '{"Group": "David Bowie", "Song": "Heroes"}'
#     '{"Group": "Fleetwood Mac", "Song": "Go Your Own Way"}'
#     '{"Group": "The Doors", "Song": "Light My Fire"}'
#     '{"Group": "R.E.M.", "Song": "Losing My Religion"}'
#     '{"Group": "Oasis", "Song": "Wonderwall"}'
#     '{"Group": "The Cure", "Song": "Lovesong"}'
#     '{"Group": "Pearl Jam", "Song": "Alive"}'
#     '{"Group": "The Smashing Pumpkins", "Song": "1979"}'
#     '{"Group": "Blur", "Song": "Song 2"}'
#     '{"Group": "Green Day", "Song": "Boulevard of Broken Dreams"}'
#     '{"Group": "Foo Fighters", "Song": "Everlong"}'
#     '{"Group": "Linkin Park", "Song": "In the End"}'
#     '{"Group": "Arctic Monkeys", "Song": "Do I Wanna Know?"}'
#     '{"Group": "Kanye West", "Song": "Stronger"}'
#     '{"Group": "Daft Punk", "Song": "Get Lucky"}'
#     '{"Group": "The Weeknd", "Song": "Blinding Lights"}'
#     '{"Group": "Billie Eilish", "Song": "Bad Guy"}'
#     '{"Group": "Lady Gaga", "Song": "Poker Face"}'
#     '{"Group": "Taylor Swift", "Song": "Shake It Off"}'
#     '{"Group": "Beyoncé", "Song": "Single Ladies"}'
#     '{"Group": "Rihanna", "Song": "Umbrella"}'
#     '{"Group": "Katy Perry", "Song": "Roar"}'
#     '{"Group": "Shakira", "Song": "Hips Dont Lie"}'
#     '{"Group": "Bruno Mars", "Song": "Uptown Funk"}'
#     '{"Group": "Ed Sheeran", "Song": "Shape of You"}'
#     '{"Group": "Maroon 5", "Song": "Sugar"}'
#     '{"Group": "Imagine Dragons", "Song": "Radioactive"}'
#     '{"Group": "OneRepublic", "Song": "Counting Stars"}'
#     '{"Group": "Post Malone", "Song": "Circles"}'
#     '{"Group": "Drake", "Song": "Gods Plan"}'
#     '{"Group": "Eminem", "Song": "Lose Yourself"}'
#     '{"Group": "Jay-Z", "Song": "Empire State of Mind"}'
#     '{"Group": "Kendrick Lamar", "Song": "HUMBLE."}'
#     '{"Group": "Travis Scott", "Song": "Sicko Mode"}'
# )

# for data in "${data_list[@]}"
# do
#     response=$(curl -s -o /dev/null -w "%{http_code}" -X POST $url -H "$headers" -d "$data")
#     echo "Status Code: $response, Data Sent: $data"
# done
# url="http://localhost:8080/api/v1/update"
# headers="Content-Type: application/json"
# curl -X PUT ${url} \
#     -H "${headers}" \
#     -d '{
#         "song_id": 48,
#         "song_detail": {
#             "release_date": "11.09.2001",
#             "link": "https://www.around_the_world.com",
#             "text": "Around The World"
#         }
#     }'
# url="http://localhost:8080/api/v1/delete"
# headers="Content-Type: application/json"
# curl -X DELETE ${url} \
#     -H "${headers}" \
#     -d '{ "song_id": 48 }'
# unit testi dlya lohov))0)
