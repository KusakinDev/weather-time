"use client";

import React, { useState, useEffect } from "react";
import { weatherType } from "@/types/weatherType";
import { ErrorType } from "@/types/errorType";
import axios from "axios";
import API_URL from "@/config";
import { FaSearch, FaTrash } from "react-icons/fa";
import Image from 'next/image';

type WeatherProps = {
  initData: weatherType;
};



const WeatherCard: React.FC<WeatherProps> = ({ initData }) => {

  const initialCity = initData.name !== " " ? initData.name : "";

  const [city, setCity] = useState(initialCity);
  const [currentTime, setCurrentTime] = useState("");
  const [errorMessage, setErrorMessage] = useState<ErrorType>({"message": "no error"})

  const [weatherType, setWeather] = useState<weatherType>(initData);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setCity(e.target.value);
  };

  const handleSearch = async () => {
    console.log("Ищем погоду для города:", city);

    try {
      const response = await axios.get(`${API_URL}/weather?name=${city}`);
      setWeather(response.data); 
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        if (error.response.status === 400) {
          const errorMessage = error.response.data.message || 'Unknown error';
          setErrorMessage(errorMessage);
          alert(errorMessage);
        } else {
          console.error('Unhandled response status:', error.response.status);
        }
      } else {
        console.error('Request error:', error);
      }
    }
  };

  const handleClear = () => {
    setCity("");
  };

  useEffect(() => {
    console.log('Данные о погоде:', weatherType);
  }, [weatherType]);
  
  useEffect(() => {
    console.log('Сообщение об ошибке:', errorMessage);
    setErrorMessage(errorMessage);
  }, [errorMessage]);

  useEffect(() => {
    const updateTime = () => {

      const now = new Date();
      const utcHours = now.getUTCHours();
      const utcMinutes = now.getUTCMinutes();
      const utcSeconds = now.getUTCSeconds();


      const adjustedTime = new Date();
      adjustedTime.setUTCHours(utcHours + weatherType.timezone);
      adjustedTime.setUTCMinutes(utcMinutes);
      adjustedTime.setUTCSeconds(utcSeconds);

      setCurrentTime(adjustedTime.toLocaleTimeString("en-GB", { timeZone: "UTC" }));
    };

    updateTime();
    const intervalId = setInterval(updateTime, 1000);

    return () => clearInterval(intervalId);
  }, [weatherType.timezone]);

  useEffect(() => {
    if (initData.name !== " ") {
      handleSearch();
    }
  }, [initData.name]);

  return (
    <div className="w-full max-w-lg h-auto mx-auto bg-bgComp shadow-md rounded-xl p-4">
      <div className="mb-4">
        <label htmlFor="city" className="block text-m text-center font-bold text-txComp">
          Узнай погоду в городе 
        </label>
        <div className="relative mt-4">
          <input
            type="text"
            id="city"
            value={city}
            onChange={handleInputChange}
            placeholder={`Введите название города...`}
            className="bg-bgElem w-full h-10 rounded-md shadow-sm sm:text-m text-txElemActive pl-3 focus:outline-none focus:border-bgElemActive focus:ring-2 focus:ring-bgElemActive"
          />
          <div className="absolute inset-y-0 right-0 flex items-center space-x-2 pr-2">
            <button
              onClick={handleSearch}
              className="p-1 text-txElem hover:text-txElemActive focus:outline-none"
            >
              <FaSearch size={15} />
            </button>
            <button
              onClick={handleClear}
              className="p-1 text-txElem hover:text-txElemActive focus:outline-none"
            >
              <FaTrash size={15} />
            </button>
          </div>
        </div>
      </div>

      <div className="text-center">
        <div className="text-5xl font-bold text-txComp">{currentTime}</div>
      </div>

      <div className="mt-4 flex justify-center gap-4">
        <div className="flex flex-col items-center pt-2 bg-bgElem rounded-md shadow text-txElem">
          <img src={weatherType.icon} alt="Weather icon" className="w-55 h-55" />
        </div>
        <div className="flex flex-col items-center pt-8 pl-4 pr-4 pb-8 bg-bgElem rounded-md shadow text-txElem">
          <div className="text-4xl font-semibold">{weatherType.temp}°C</div>
        </div>
      </div>
    </div>
  );
};

export default WeatherCard;