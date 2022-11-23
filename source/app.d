/+ 
	Canarium: A simple weather app made in the D programming language. 
	The appearance of the app takes inspiration from the WeatherStar 3000 
	and WeatherStar JR systems used by The Weather Channel.

	System Requirements:
		Windows and Linux:
		192 MB of RAM (for running and building the app, along with X11/Wayland).
		512 MHz CPU, anything from the past 20-25 years should suffice.

	Windows XP and higher is supported.
+/

import std.file;
import std.json;
import std.conv;
import std.stdio;
import std.ascii;
import std.string;
import std.format;
import std.random;
import std.exception;
import std.algorithm;
import std.datetime;

import arsd.simpledisplay;
import arsd.simpleaudio;
import arsd.http2;
import arsd.mp3;
import arsd.rss;
import arsd.dom;
import arsd.ttf;

JSONValue settings;

JSONValue downloadParseData(int mode = 0)
{
	settings = parseJSON(readText("config.json"));

	final switch (mode)
	{
	case 0:
		auto client = new HttpClient();
		auto request = client.request(Uri(format(
				"https://api.weather.com/v1/location/%s:4:US/observations/current.json?language=en-US&units=e&apiKey=%s",
				settings["zip"].str, settings["wkey"].str)));
		auto response = request.waitForCompletion();
		return parseJSON(response.contentText);

	case 1:
		auto client = new HttpClient();
		auto request = client.request(Uri(format(
				"https://nominatim.openstreetmap.org/search?city=%s&state%s&postalcode=%s&format=json",
				settings["city"].str, settings["state"].str, settings["zip"].str)));
		auto response = request.waitForCompletion();
		JSONValue geo = parseJSON(response.contentText);

		request = client.request(Uri(format(
				"https://api.weather.com/v3/aggcommon/v3-wx-forecast-daily-5day?geocodes=%s,%s&language=en-US&units=e&format=json&apiKey=%s", 
				geo[0]["lat"].str, geo[0]["lon"].str, settings["wkey"].str)));
		response = request.waitForCompletion();

		return parseJSON(response.contentText);
	}
}

RssChannel parseRSSFeed()
{
	settings = parseJSON(readText("config.json"));
	auto client = new HttpClient();
	auto request = client.request(Uri(settings["feed"].str));
	auto response = request.waitForCompletion();
	if (response.code != 200)
	{
		final switch (response.code)
		{
		case 301:
		case 302:
			throw new Exception(
				"Redirect. A common mistake would usually be your RSS URL using the HTTP protocol and not HTTPS. 
				This may also mean your RSS isn't valid.");
		case 400:
		case 404:
			throw new Exception(
				"Client Error. Either your RSS URL doesn't exist, or is malformed.");
		case 408:
			throw new Exception(
				"RSS Feed is probably down.");
		case 500:
			throw new Exception(
				"The URL for the RSS feed is currently experiencing some sort of server error.");
		}
	}

	else
	{
		return parseRss(to!string(response.contentText));
	}
}

@safe string clock()
{
	SysTime today = Clock.currTime();
	return format("%02d:%02d:%02d %s", today.hour > 12 ? today.hour - 12 : today.hour, today.minute, today.second, today
			.hour >= 12 ? "PM" : "AM");
}

void main()
{
	static DrawableFont StarJR, StarJRNarrative, StarJRHead, StarJRSS;

	if (StarJRHead is null || StarJRNarrative is null || StarJR is null)
	{
		StarJR = arsdTtfFont(cast(ubyte[]) std.file.read("assets/fonts/StarJR.ttf"), 34);
		StarJRNarrative = arsdTtfFont(cast(ubyte[]) std.file.read("assets/fonts/StarJR.ttf"), 28);
		StarJRHead = arsdTtfFont(cast(ubyte[]) std.file.read("assets/fonts/StarJR.ttf"), 58);
		StarJRSS = arsdTtfFont(cast(ubyte[]) std.file.read("assets/fonts/StarJR.ttf"), 22);
	}

	string[] musicbox()
	{
		string[] music;

		foreach (string name; dirEntries("assets\\music", SpanMode.breadth))
		{
			music ~= name;
		}

		auto rnd0 = MinstdRand0(1);
		rnd0.seed(unpredictableSeed);

		music = music.randomShuffle(rnd0);
		return music;
	}

	auto audio = AudioOutputThread(true);
	string[] list = musicbox();
	SampleController sc = audio.playMp3(list[0]);

	void Draw(SimpleWindow window)
	{
		auto painter = window.draw();

		JSONValue wxData = downloadParseData();
		JSONValue narrative = downloadParseData(1);

		RssItem[] news = parseRSSFeed().items;
		settings = parseJSON(readText("config.json"));

		painter.clear(Color(0, 0, 120));
		painter.outlineColor = Color.white;
		auto point = Point(16, 16);

		void header(string txt)
		{
			painter.fillColor = Color(0, 0, 120);
			painter.outlineColor = Color.white;

			painter.drawText(StarJRHead, point, txt);
			point.y += StarJRHead.height();
			point.y += 20;
		}

		void text(string txt)
		{
			painter.fillColor = Color(0, 0, 120);
			painter.outlineColor = Color.white;

			painter.drawText(StarJR, point, txt);
			point.y += StarJR.height();
			point.y += 10;
		}

		void outlook(string txt)
		{
			painter.fillColor = Color(0, 0, 120);
			painter.outlineColor = Color.white;

			painter.drawText(StarJRNarrative, point, txt);
			point.y += StarJRNarrative.height();
			point.y += 15;
		}

		void rss(string txt)
		{
			painter.fillColor = Color(0, 0, 120);
			painter.outlineColor = Color.white;

			painter.drawText(StarJRSS, point, txt);
			point.y += StarJRNarrative.height();
		}

		header("Current Conditions");
		text(format("Temperature: %s%sF", wxData["observation"]["imperial"]["hi"], "°"));
		text(format("Dew Point: %s%sF", wxData["observation"]["imperial"]["dewpt"], "°"));
		text(format("Wind Speed: %smph", wxData["observation"]["imperial"]["wspd"]));
		text(format("Visibility: %s miles", to!string(
				wxData["observation"]["imperial"]["vis"].floating)));
		text(format("UV Index: %s (%s)", wxData["observation"]["uv_index"], wxData["observation"]["uv_desc"]
				.str));
		text(format("Cloud Cover: %s", wxData["observation"]["sky_cover"].str));
		outlook(format("Today's Forecast: %s", narrative[0]["v3-wx-forecast-daily-5day"]["narrative"][0]
				.str));
		header("Local Channels");
		text(format("%s %s", settings["channels"][0]["id"].str, settings["channels"][0]["name"]
				.str));
		text(format("%s %s", settings["channels"][1]["id"].str, settings["channels"][1]["name"]
				.str));
		text(format("%s %s", settings["channels"][2]["id"].str, settings["channels"][2]["name"]
				.str));
		text(format("%s %s", settings["channels"][3]["id"].str, settings["channels"][3]["name"]
				.str));
		header("Today's Top News Article");
		text(format("%s", news[0].title));

		foreach (line; splitLines(wrap(news[0].description, 100)))
		{
			rss(htmlEntitiesDecode(to!string(line)));
		}

		writeln("Updated!");
	}

	void hms(SimpleWindow window)
	{
		auto painter = window.draw();
		auto point = Point(475, 35);

		void text(string txt)
		{
			painter.fillColor = Color(0, 0, 120);
			painter.outlineColor = Color.white;

			painter.drawText(StarJR, point, txt);
		}

		text(clock());
	}

	void music()
	{
		if (sc.finished)
		{
			if (list == null)
			{
				list = musicbox();
			}
			else
			{
				list = list.remove(0);
			}

			sc = audio.playMp3(list[0]);
		}
	}

	auto window = new SimpleWindow(960, 720, "Canarium");
	auto timer = new Timer(600_000, delegate{ Draw(window); });
	auto time = new Timer(10_00, delegate{ hms(window); });
	auto musicheck = new Timer(10_00, delegate{ music(); });

	window.maximize();
	Draw(window);

	window.eventLoop(0);

	scope (exit)
	{
		sc.stop();
	}
}
