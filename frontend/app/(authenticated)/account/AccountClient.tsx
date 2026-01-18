import ChangePasswordSection from "./ChangePasswordSection";
import DeleteAccountSection from "./DeleteAccountSection";
import ThemeSection from "./ThemeSection";
import UserProfileSection from "./UserProfileSection";

export default function AccountClient() {
  return (
    <>
      <h1 className="sr-only">Account Settings</h1>

      {/* Settings forms */}
      <div className="divide-y divide-white/5">
        <UserProfileSection />
        <ThemeSection />
        <ChangePasswordSection />
        <DeleteAccountSection />
      </div>
    </>
  )
}
