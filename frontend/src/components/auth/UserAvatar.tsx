'use client';

import { useState, useRef, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/lib/auth/auth-context';
import Image from 'next/image';

export function UserAvatar() {
    const { user, logout } = useAuth();
    const [open, setOpen] = useState(false);
    const dropdownRef = useRef<HTMLDivElement>(null);
    const router = useRouter();

    // Close dropdown when clicking outside of it
    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
                setOpen(false);
            }
        };

        window.addEventListener('click', handleClickOutside);
        return () => {
            window.removeEventListener('click', handleClickOutside);
        };
    }, []);

    if (!user) return null;

    return (
        <div className="relative inline-flex items-center group" ref={dropdownRef}>
            <button
                onClick={() => setOpen(!open)}
                className="flex items-center h-10 focus:outline-none -mt-8"
            >
                <Image
                    src={user.avatarUrl || '/default-avatar.png'}
                    alt="User Avatar"
                    width={40}
                    height={40}
                    className="rounded-full"
                />
            </button>
            {/* Tooltip visible on hover */}
            <span className="absolute bottom-0 left-1/2 transform -translate-x-1/2 translate-y-full opacity-0 group-hover:opacity-100 transition bg-gray-800 text-white text-xs px-2 py-1 rounded mt-1 whitespace-nowrap z-10">
                {user.username}
            </span>
            {open && (
                <div className="absolute right-0 mt-20 w-40 bg-white border border-gray-200 rounded shadow-lg">
                    <button
                        onClick={async () => {
                            setOpen(false);
                            await logout();
                            router.push('/'); // redirect to the home page after logout
                        }}
                        className="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                    >
                        Logout
                    </button>
                </div>
            )}
        </div>
    );
}