//
//  WHOMessage.m
//  Hooli
//
//  Created by dylan on 4/19/14.
//  Copyright (c) 2014 whoisdylan. All rights reserved.
//

#import "WHOMessage.h"

@implementation WHOMessage

-(BOOL)isEqual:(id)object {
    WHOMessage* anObject = (WHOMessage*) object;
    if (![self.author isEqualToString:anObject.author]) {
        NSLog(@"not equal 1");
        return NO;
    }
    else if (![self.message isEqualToString:anObject.message]) {
        NSLog(@"not equal 2");
        return NO;
    }
    double distanceApartInMiles = ([self.location distanceFromLocation:anObject.location])/1609.344;
    if (distanceApartInMiles > 0.1) {
        NSLog(@"not equal 3");
        return NO;
    }
    else {
//        NSLog(@"message equal");
        return YES;
    }
}

- (instancetype)initWithMessage:(NSString *)message Author:(NSString *)author Location:(CLLocation *)location {
    if (self = [super init]) {
        self.message = message;
        self.author = author;
        self.location = location;
    }
    return self;
}

@end
