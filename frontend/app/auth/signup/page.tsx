import { Metadata } from "next";
import SignupClient from "./SignupClient";

export const metadata: Metadata = {
  title: "Sign Up",
  description: "Create your TuneSlap account",
};

export default function SignupPage() {
  return <SignupClient />;
}
