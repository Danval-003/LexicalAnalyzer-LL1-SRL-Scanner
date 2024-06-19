import axios from 'axios';
import { API_URL } from '../env';  // Asegúrate de que la ruta de importación es correcta

console.log("Api", API_URL);

export const CreateScanner = async (yalex: string) => {
  try {
    const response = await axios.post(`${API_URL}/scanners/public/create`, {
      content: yalex
    }, {
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      }
    });
    console.log("Response data:", response.data);
    return response.data;
  } catch (error) {
    console.error("Error creating scanner:", error);
    throw error;
  }
};

export const CreateSLR = async (yapar: string) => {
  // Use /yapar/pub/create
  try {
    const response = await axios.post(`${API_URL}/yapar/pub/create`, {
      content: yapar
    }, {
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      }
    });
    console.log("Response data:", response.data);
    return response.data;
  } catch (error) {
    console.error("Error creating SLR:", error);
    throw error;
  }
}

export const SimulateText = async (text: string, scannerName: string, SLR: string) => {
  const body = {
    contentSimulate: text,
    scannerName: scannerName,
    slrName: SLR
  };

  console.log("Body:", body);
  try {
    const response = await axios.post(`${API_URL}/compare/simulate`, body, {
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      }
    });
    console.log("Response data:", response.data);
    return response.data;
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      const errorMessage = error.response.data.message;
      const errorStatus = error.response.data.status;
      console.error(`Error ${errorStatus}: ${errorMessage}`);
      throw new Error(errorMessage);
    } else {
      console.error("Unknown error:", error);
      throw new Error("Unknown error occurred");
    }
  }
};
